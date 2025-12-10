using System;

namespace Tachograph
{
    /// <summary>
    /// Configures the anonymization process for tachograph files.
    /// </summary>
    public class AnonymizeOptions
    {
        /// <summary>
        /// Controls whether distance and trip data are preserved.
        /// 
        /// If true, odometer readings and distance values are preserved in their original form.
        /// If false (default), distance data is rounded or anonymized to obscure exact values.
        /// </summary>
        public bool PreserveDistanceAndTrips { get; set; } = false;

        /// <summary>
        /// Controls whether timestamps are preserved.
        /// 
        /// If true, timestamps are preserved in their original form.
        /// If false (default), timestamps are shifted to a fixed epoch (2020-01-01 00:00:00 UTC)
        /// to obscure the exact time of events while maintaining relative ordering.
        /// </summary>
        public bool PreserveTimestamps { get; set; } = false;

        /// <summary>
        /// Creates an anonymized copy of a parsed tachograph file.
        /// 
        /// Anonymization replaces personally identifiable information (PII) with test values
        /// while preserving the structural integrity of the file for testing purposes.
        /// 
        /// The zero value of AnonymizeOptions anonymizes both timestamps and distances.
        /// </summary>
        /// <param name="file">The file to anonymize</param>
        /// <returns>An anonymized File object</returns>
        /// <exception cref="ArgumentNullException">Thrown when file is null</exception>
        /// <exception cref="NotSupportedException">Thrown when the file type is not supported</exception>
        public File Anonymize(File file)
        {
            if (file == null)
                throw new ArgumentNullException(nameof(file));

            var result = new File
            {
                Type = file.Type
            };

            switch (file.Type)
            {
                case FileType.DriverCard:
                    result.DriverCard = AnonymizeDriverCardFile(file.DriverCard);
                    break;

                case FileType.VehicleUnit:
                    result.VehicleUnit = AnonymizeVehicleUnitFile(file.VehicleUnit);
                    break;

                default:
                    throw new NotSupportedException($"Unsupported file type for anonymization: {file.Type}");
            }

            return result;
        }

        private DriverCardFile AnonymizeDriverCardFile(DriverCardFile card)
        {
            // TODO: Implement driver card file anonymization
            throw new NotImplementedException("Driver card file anonymization not yet implemented");
        }

        private VehicleUnitFile AnonymizeVehicleUnitFile(VehicleUnitFile vu)
        {
            // TODO: Implement VU file anonymization
            throw new NotImplementedException("Vehicle unit file anonymization not yet implemented");
        }
    }
}
