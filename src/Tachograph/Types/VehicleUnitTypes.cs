using System;

namespace Tachograph
{
    /// <summary>
    /// Represents a raw vehicle unit file.
    /// </summary>
    public class RawVehicleUnitFile
    {
        /// <summary>
        /// Gets or sets the raw binary data of the vehicle unit file.
        /// </summary>
        public byte[] Data { get; set; }

        /// <summary>
        /// Creates a deep clone of this RawVehicleUnitFile.
        /// </summary>
        /// <returns>A cloned RawVehicleUnitFile</returns>
        public RawVehicleUnitFile Clone()
        {
            return new RawVehicleUnitFile
            {
                Data = (byte[])this.Data?.Clone()
            };
        }
    }

    /// <summary>
    /// Represents a parsed vehicle unit file.
    /// </summary>
    public class VehicleUnitFile
    {
        // TODO: Add vehicle unit specific fields
    }
}
